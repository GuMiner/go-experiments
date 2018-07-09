using System;
using System.Drawing;
using System.Drawing.Imaging;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class Form2 : Form, IDisposable
    {
        private const int scaleDownFactor = 2;

        private readonly int subImageId;
        private readonly SubImage subImage;
        private ImageAttributes imageAttributes;

        private bool inTranslateMode = false;
        private Point lastMousePos = new Point(-1, -1);

        private ImageEditForm imageEditForm;

        // 0 == no rotation, 1-3, rotation, 
        private int rotation = 0;

        public Form2()
        {
            InitializeComponent();
            this.imageAttributes = new ImageAttributes();
            this.MouseWheel += new MouseEventHandler(MouseWheelMove);
        }

        private void MouseWheelMove(object sender, MouseEventArgs e)
        {
            float aspectRatio = (float)(subImage.MaxPos.X - subImage.MinPos.X) / (float)(subImage.MaxPos.Y - subImage.MinPos.Y);

            int scrollDelta = e.Delta;
            int zoomAddFactor = 10;
            if (scrollDelta > 0)
            {
                subImage.MinPos = new Point((int)((float)subImage.MinPos.X - zoomAddFactor * aspectRatio), (int)((float)subImage.MinPos.Y - zoomAddFactor));
                subImage.MaxPos = new Point((int)((float)subImage.MaxPos.X + zoomAddFactor * aspectRatio), (int)((float)subImage.MaxPos.Y + zoomAddFactor));
            }
            else
            {
                subImage.MinPos = new Point((int)((float)subImage.MinPos.X + zoomAddFactor * aspectRatio), (int)((float)subImage.MinPos.Y + zoomAddFactor));
                subImage.MaxPos = new Point((int)((float)subImage.MaxPos.X - zoomAddFactor * aspectRatio), (int)((float)subImage.MaxPos.Y - zoomAddFactor));
            }

            imageBox.Invalidate();
        }

        public Form2(int subImageId, SubImage subImage)
            : this()
        {
            this.Text = $"Sub Image - {subImageId}";
            this.subImageId = subImageId;
            this.subImage = subImage;

            this.imageEditForm = new ImageEditForm(
                this.subImageId,
                (attr) =>
                {
                    this.imageAttributes = attr;
                    imageBox.Invalidate();
                },
                () => {
                    int width = subImage.MaxPos.X - subImage.MinPos.X;
                    int height = subImage.MaxPos.Y - subImage.MinPos.Y;
                    if (rotation % 4 == 1 || rotation % 4 == 3)
                    {
                        int swap = width;
                        width = height;
                        height = swap;
                    }

                    Bitmap bitmap = new Bitmap(width, height);
                    using (Graphics g = Graphics.FromImage(bitmap))
                    {
                        DrawPartialImage(g, width, height);
                    }

                    this.subImage.Saved = true;
                    return bitmap;
                },
                (increment) =>
                {
                    this.rotation += increment;
                    this.Width = this.GetWidth();
                    this.Height = this.GetHeight();

                    imageBox.Invalidate();
                },
                () => this.Hide());

            this.Width = this.GetWidth();
            this.Height = this.GetHeight();
            this.Show();

            this.imageEditForm.StartPosition = FormStartPosition.Manual;
            this.imageEditForm.Left = this.Location.X + this.Width;
            this.imageEditForm.Top = this.Location.Y;
            this.imageEditForm.Show();
        }

        public new void Dispose()
        {
            this.imageBox.Image = null;
            this.imageEditForm.Dispose();
            base.Dispose();
        }

        private void pictureBox1_MouseMove(object sender, MouseEventArgs e)
        {
            if (inTranslateMode)
            {
                int deltaX = e.X - lastMousePos.X;
                int deltaY = e.Y - lastMousePos.Y;
                switch (rotation % 4)
                {
                    case 0:
                        break;
                    case 1:
                        int swap = deltaX;
                        deltaX = -deltaY;
                        deltaY = swap;
                        break;
                    case 2:
                        deltaX = -deltaX;
                        deltaY = -deltaY;
                        break;
                    case 3:
                        swap = deltaX;
                        deltaX = deltaY;
                        deltaY = -swap;
                        break;
                    default:
                        break;
                }

                float xScaleFactor = (float)(subImage.MaxPos.X - subImage.MinPos.X) / (float)imageBox.Width;
                float yScaleFactor = (float)(subImage.MaxPos.Y - subImage.MinPos.Y) / (float)imageBox.Height;
                subImage.MinPos = new Point((int)((float)subImage.MinPos.X - (float)deltaX * xScaleFactor), (int)((float)subImage.MinPos.Y - (float)deltaY * yScaleFactor));
                subImage.MaxPos = new Point((int)((float)subImage.MaxPos.X - (float)deltaX * xScaleFactor), (int)((float)subImage.MaxPos.Y - (float)deltaY * yScaleFactor));

                lastMousePos = e.Location;
                imageBox.Invalidate();
            }
        }

        private void pictureBox1_MouseDown(object sender, MouseEventArgs e)
        {
            inTranslateMode = true;
            lastMousePos = e.Location;
        }

        private void pictureBox1_MouseUp(object sender, MouseEventArgs e)
        {
            inTranslateMode = false;
        }

        private void imageBox_Paint(object sender, PaintEventArgs e)
        {
            DrawPartialImage(e.Graphics, imageBox.Width, imageBox.Height);
        }

        private void DrawPartialImage(Graphics g, int destinationWidth, int destinationHeight)
        {
            g.Clear(Color.LightSeaGreen);
            g.DrawImage(subImage.Image,
                this.GetDestinationPoints(destinationWidth, destinationHeight),
                new Rectangle(subImage.MinPos.X, subImage.MinPos.Y, subImage.MaxPos.X - subImage.MinPos.X, subImage.MaxPos.Y - subImage.MinPos.Y),
                GraphicsUnit.Pixel,
                imageAttributes);
        }

        private void Form2_Resize(object sender, EventArgs e)
        {
            imageBox.Invalidate();
        }

        private PointF[] GetDestinationPoints(int destinationWidth, int destinationHeight)
        {
            --destinationWidth;
            --destinationHeight;
            switch (rotation % 4)
            {
                case 0:
                    return new[] { new PointF(0, 0), new PointF(destinationWidth, 0), new PointF(0, destinationHeight) };
                case 1:
                    return new[] { new PointF(destinationWidth, 0), new PointF(destinationWidth, destinationHeight), new PointF(0, 0),  };
                case 2:
                    return new[] { new PointF(destinationWidth, destinationHeight), new PointF(0, destinationHeight), new PointF(destinationWidth, 0), };
                case 3:
                default:
                    return new[] { new PointF(0, destinationHeight), new PointF(0, 0), new PointF(destinationWidth, destinationHeight), };
            }
        }

        private int GetWidth()
        {
            int nominalWidth = subImage.MaxPos.X - subImage.MinPos.X;
            int nominalHeight = subImage.MaxPos.Y - subImage.MinPos.Y;
            switch (rotation % 4)
            {
                case 0:
                case 2:
                    return nominalWidth / scaleDownFactor;
                case 1:
                case 3:
                default:
                    return nominalHeight / scaleDownFactor;
            }
        }

        private int GetHeight()
        {
            int nominalWidth = subImage.MaxPos.X - subImage.MinPos.X;
            int nominalHeight = subImage.MaxPos.Y - subImage.MinPos.Y;
            switch (rotation % 4)
            {
                case 0:
                case 2:
                    return nominalHeight / scaleDownFactor;
                case 1:
                case 3:
                default:
                    return nominalWidth / scaleDownFactor;
            }
        }
    }
}
