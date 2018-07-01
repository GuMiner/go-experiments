using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Drawing.Imaging;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class Form2 : Form
    {
        private SubImage subImage;
        private ImageAttributes imageAttributes;
        public Form2()
        {
            InitializeComponent();
            this.imageAttributes = new ImageAttributes();
            this.MouseWheel += new MouseEventHandler(MouseWheelMove);

            ImageEditForm imageEditForm = new ImageEditForm((attr) =>
            {
                this.imageAttributes = attr;
                imageBox.Invalidate();
            });
            imageEditForm.Show();
        }

        private void MouseWheelMove(object sender, MouseEventArgs e)
        {
            float aspectRatio = 1.0f / ((float)subImage.Image.Width / (float)subImage.Image.Height);
            int scrollDelta = e.Delta;
            int zoomAddFactor = 10;
            if (scrollDelta > 0)
            {
                subImage.MinPos = new Point((int)((float)subImage.MinPos.X - zoomAddFactor * aspectRatio), (int)((float)subImage.MinPos.Y - zoomAddFactor));
                subImage.MaxPos = new Point((int)((float)subImage.MaxPos.X + zoomAddFactor * aspectRatio), (int)((float)subImage.MaxPos.Y + zoomAddFactor));
            }
            else
            {
                subImage.MinPos = new Point((int)((float)subImage.MinPos.X + zoomAddFactor* aspectRatio), (int)((float)subImage.MinPos.Y + zoomAddFactor));
                subImage.MaxPos = new Point((int)((float)subImage.MaxPos.X - zoomAddFactor* aspectRatio), (int)((float)subImage.MaxPos.Y - zoomAddFactor));
            }

            imageBox.Invalidate();
        }

        public Form2(SubImage subImage)
            : this()
        {
            this.subImage = subImage;
        }

        private void pictureBox1_MouseMove(object sender, MouseEventArgs e)
        {

        }

        private void pictureBox1_MouseDown(object sender, MouseEventArgs e)
        {

        }

        private void pictureBox1_MouseUp(object sender, MouseEventArgs e)
        {

        }

        private void imageBox_Paint(object sender, PaintEventArgs e)
        {
            e.Graphics.Clear(Color.LightSeaGreen);
            e.Graphics.DrawImage(subImage.Image,
                new Rectangle(0, 0, imageBox.Width, imageBox.Height),
                subImage.MinPos.X, subImage.MinPos.Y,
                subImage.MaxPos.X - subImage.MinPos.X,
                subImage.MaxPos.Y - subImage.MinPos.Y,
                GraphicsUnit.Pixel,
                imageAttributes);
        }

        private void Form2_Resize(object sender, EventArgs e)
        {
            imageBox.Invalidate();
        }
    }
}
